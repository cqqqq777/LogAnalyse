package main

import (
	"context"
	"errors"
	"time"

	"LogAnalyse/app/service/user/internal"
	"LogAnalyse/app/service/user/rpc/dao"
	"LogAnalyse/app/service/user/rpc/model"
	"LogAnalyse/app/service/user/rpc/pkg"
	"LogAnalyse/app/shared/consts"
	"LogAnalyse/app/shared/errz"
	"LogAnalyse/app/shared/kitex_gen/user"
	"LogAnalyse/app/shared/log"
	"LogAnalyse/app/shared/middleware"
	"LogAnalyse/app/shared/response"

	"github.com/bwmarrin/snowflake"
	"github.com/golang-jwt/jwt"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	Dao *dao.User
	JWT *middleware.JWT
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	// TODO: Your code here...
	resp = new(user.RegisterResp)

	sf, err := snowflake.NewNode(internal.UserSnowflakeNode)
	if err != nil {
		log.Zlogger.Error("generate user id failed err:" + err.Error())
		resp.BaseResp = response.NewBaseResp(errz.ErrGenerateID)
		return resp, nil
	}

	id := sf.Generate().Int64()
	User := &model.User{
		Id:       id,
		Username: req.Username,
		Password: pkg.Md5(req.Password),
	}

	// create user in mysql
	if err = s.Dao.CreateUser(User); err != nil {
		if errors.Is(err, dao.ErrUserExist) {
			resp.BaseResp = response.NewBaseResp(errz.ErrUserExist)
			return resp, nil
		} else {
			log.Zlogger.Error("create user failed err:" + err.Error())
			resp.BaseResp = response.NewBaseResp(errz.ErrUserService)
			return resp, nil
		}
	}

	token, err := s.JWT.CreateToken(middleware.CustomClaims{
		ID:       id,
		Identity: consts.UserIdentity,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + consts.TokenExpiredAt,
		},
	})
	if err != nil {
		log.Zlogger.Error("generate token failed err:" + err.Error())
		resp.BaseResp = response.NewBaseResp(errz.ErrGenerateToken)
		return resp, nil
	}

	resp.BaseResp = response.NewBaseResp(nil)
	resp.Token = token
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	// TODO: Your code here...
	resp = new(user.LoginResp)

	User, err := s.Dao.GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, dao.ErrUserNotFound) {
			resp.BaseResp = response.NewBaseResp(errz.ErrUserNotFound)
			return resp, nil
		}
		resp.BaseResp = response.NewBaseResp(errz.ErrUserService)
		log.Zlogger.Error("get user by username failed err:" + err.Error())
		return resp, nil
	}

	if User.Password != pkg.Md5(req.Password) {
		resp.BaseResp = response.NewBaseResp(errz.ErrWrongPassword)
		return resp, nil
	}

	token, err := s.JWT.CreateToken(middleware.CustomClaims{
		ID:       User.Id,
		Identity: consts.UserIdentity,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + consts.TokenExpiredAt,
		},
	})
	if err != nil {
		log.Zlogger.Error("generate token failed err:" + err.Error())
		resp.BaseResp = response.NewBaseResp(errz.ErrGenerateToken)
		return resp, nil
	}

	resp.BaseResp = response.NewBaseResp(nil)
	resp.Token = token
	resp.Id = User.Id
	return
}
